use tonic::{transport::Server, Request, Response, Status};
extern crate title_fetcher;

use std::io;
use tonic::{transport::Server, Code, Request, Response, Status};
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
    async fn fetch(
        &self,
        request: Request<FetchRequest>,
    ) -> Result<Response<FetchReply>, Status> {
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
        };
        Ok(Response::new(reply))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse().unwrap();
    let title_fetcher = MyTitleFetcher::default();

    println!("TitleFetcherServer listening on {}", addr);

    Server::builder()
        .add_service(TitleFetcherServer::new(title_fetcher))
        .serve(addr)
        .await?;

    Ok(())
}
