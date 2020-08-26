extern crate title_fetcher;

use std::io;
use std::time::Duration;
use tokio::time::delay_for;
use tonic::{transport::Server, Code, Request, Response, Status};
use tonic_health::server::HealthReporter;

use pb::title_fetcher_server::{TitleFetcher, TitleFetcherServer};
use pb::{FetchReply, FetchRequest};

pub mod pb {
    tonic::include_proto!("title_fetcher");
}

#[derive(Default)]
pub struct MyTitleFetcher {}

enum Error {
    HTTP(reqwest::StatusCode),
    Internal,
    FailedToSerialize,
}

// TODO test fetch_title
// how to test async function?
async fn fetch_title(url: &str) -> Result<String, Error> {
    let res = reqwest::get(url)
        .await
        .map_err(|e| e.status().map_or_else(|| Error::Internal, Error::HTTP))?;
    let body = res.text().await.map_err(|_| Error::FailedToSerialize)?;
    let title = title_fetcher::parser::parse(&mut io::Cursor::new(body));
    Ok(title.unwrap_or_else(String::new))
}

#[tonic::async_trait]
impl TitleFetcher for MyTitleFetcher {
    async fn fetch(&self, request: Request<FetchRequest>) -> Result<Response<FetchReply>, Status> {
        println!("Got a request from {:?}", request.remote_addr());
        match fetch_title(&request.into_inner().url).await {
            Ok(title) => Ok(Response::new(pb::FetchReply { title })),
            Err(Error::HTTP(status)) => Err(Status::new(
                Code::InvalidArgument,
                format!("failed to request via HTTP: {:?}", status),
            )),
            Err(Error::Internal) => Err(Status::new(Code::InvalidArgument, "Invalid argument")),
            Err(Error::FailedToSerialize) => {
                Err(Status::new(Code::InvalidArgument, "Internal Error"))
            }
        }
    }
}

// TODO test health check service
async fn twiddle_service_status(mut reporter: HealthReporter) {
    let mut iter = 0u64;
    loop {
        iter += 1;
        delay_for(Duration::from_secs(1)).await;

        if iter % 2 == 0 {
            reporter
                .set_serving::<TitleFetcherServer<MyTitleFetcher>>()
                .await;
        } else {
            reporter
                .set_not_serving::<TitleFetcherServer<MyTitleFetcher>>()
                .await;
        };
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let (mut health_reporter, health_service) = tonic_health::server::health_reporter();
    health_reporter
        .set_serving::<TitleFetcherServer<MyTitleFetcher>>()
        .await;

    tokio::spawn(twiddle_service_status(health_reporter.clone()));
    let addr = "[::1]:50051".parse().unwrap();
    let title_fetcher = MyTitleFetcher::default();

    println!("TitleFetcherServer listening on {}", addr);

    Server::builder()
        .add_service(TitleFetcherServer::new(title_fetcher))
        .add_service(health_service)
        .serve(addr)
        .await?;

    Ok(())
}
