// Copyright (C) 2023 Parseable, Inc.
// 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
// 
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

use lambda_extension::*;
use tracing::info;

async fn telemetry_extension(events: Vec<LambdaTelemetry>) -> Result<(), Error> {
    for event in events {
        match event.record {
            LambdaTelemetryRecord::Function(record) => {
                info!(telemetry_type = "function", record = ?record, "received function telemetry");
            }
            _ignore_other => {},
        }
    }

    Ok(())
}

#[tokio::main]
async fn main() -> Result<(), Error> {
    // The runtime logging can be enabled here by initializing `tracing` with `tracing-subscriber`
    // While `tracing` is used internally, `log` can be used as well if preferred.
    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::INFO)
        // disabling time is handy because CloudWatch will add the ingestion time.
        .without_time()
        .init();

    let telemetry_processor = SharedService::new(service_fn(telemetry_extension));
    
    Extension::new()
        .with_telemetry_processor(telemetry_processor)
        .run()
        .await
}