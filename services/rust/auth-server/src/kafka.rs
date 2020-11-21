use rdkafka::{ producer::{ FutureProducer, FutureRecord } };
use tokio::sync::mpsc;
use tokio::sync::mpsc::{ Receiver, Sender };
use std::time::Duration;

#[derive(Clone)]
pub struct Kafka {
    tx: Sender<String>
}

impl Kafka {
    pub async fn send(&mut self, message: String) -> anyhow::Result<()> {
        match self.tx.send(message).await {
            Ok(_) => {
                log::info!("Sent message");
                Ok(())
            },
            Err(e) => {
                log::error!("Failed to send");
                Err(anyhow::anyhow!("Failed to send message: {:?}", e))
            }
        }
    }
}

pub fn new() -> Kafka {
    let (tx, rx) = mpsc::channel(100);

    let kafka = Kafka{
        tx: tx
    };

    start(rx);

    return kafka;
}

fn start(mut rx: Receiver<String>) {
    tokio::spawn(async move {
        let producer: FutureProducer = rdkafka::ClientConfig::new()
            .set("bootstrap.servers", "pongle-cluster-kafka-bootstrap")
            .create()
            .expect("Failed to connect to kafka");

        log::info!("Running the kafka thread");

        loop {
            match rx.recv().await {
                Some(message) => {
                    let result = producer.send(
                        FutureRecord::to("my.topic")
                            .payload(&format!("{}", message))
                            .key("Key: 1")
                        , Duration::from_secs(0)
                    ).await;

                    match result {
                        Ok((val1, val2)) => {
                            log::info!("Success! {} {}", val1, val2);
                        },
                        Err((err, _)) => {
                            log::error!("That failed: {:?}", err);
                        }
                    }
                },
                None => {
                    log::error!("Empty send");
                }
            }
        }
    });
}