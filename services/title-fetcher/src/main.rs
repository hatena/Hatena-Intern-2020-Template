use tonic::{transport::Server, Request, Response, Status};

use pb::title_fetcher_server::{TitleFetcher, TitleFetcherServer};
use pb::{FetchReply, FetchRequest};

pub mod pb {
    tonic::include_proto!("title_fetcher");
}

#[derive(Default)]
pub struct MyTitleFetcher {}

#[tonic::async_trait]
impl TitleFetcher for MyTitleFetcher {
    async fn fetch(
        &self,
        request: Request<FetchRequest>,
    ) -> Result<Response<FetchReply>, Status> {
        println!("Got a request from {:?}", request.remote_addr());

        let reply = pb::FetchReply {
            title: String::from("Hello World!"),
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
