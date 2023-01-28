use clap::Parser;
use lambda_extension::{service_fn, Error, Extension, LambdaLog, LambdaLogRecord, SharedService};
use once_cell::sync::Lazy;

static CONFIG: Lazy<Config> = Lazy::new(|| Config::parse());

#[derive(Parser)]
struct Config {
    #[arg(env = "PARSEABLE_URL")]
    url: String,
    #[arg(env = "PARSEABLE_USERNAME")]
    username: String,
    #[arg(env = "PARSEABLE_PASSWORD")]
    password: String,
    #[arg(env = "PARSEABLE_STREAM")]
    stream: Option<String>,
}

async fn handler(logs: Vec<LambdaLog>) -> Result<(), Error> {
    for log in logs {
        match log.record {
            LambdaLogRecord::Function(_record) => {
                // do something with the function log record
                todo!()
            }
            LambdaLogRecord::Extension(_record) => {
                // do something with the extension log record
                todo!()
            }
            _ => (),
        }
    }

    Ok(())
}

#[tokio::main]
async fn main() -> Result<(), Error> {
    let logs_processor = SharedService::new(service_fn(handler));

    Extension::new()
        .with_logs_processor(logs_processor)
        .run()
        .await?;

    Ok(())
}
