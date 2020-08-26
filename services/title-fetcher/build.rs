fn main() {
    tonic_build::compile_protos("./pb/title_fetcher.proto").unwrap();
}
