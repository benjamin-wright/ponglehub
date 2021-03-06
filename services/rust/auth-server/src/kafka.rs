use rdkafka::{ producer::{ FutureProducer, FutureRecord } };
use tokio::sync::mpsc;
use tokio::sync::mpsc::{ Receiver, Sender };
use std::time::Duration;

#[derive(Clone)]
pub struct Kafka {
    pub tx: Sender<KafkaMessage>
}

#[derive(Debug)]
pub struct KafkaMessage {
    pub topic: String,
    pub message: String
}

#[macro_export]
macro_rules! send_to_kafka {
    ($kafka:expr, $topic:expr, $message:expr) => {
        if let Err(e) = $kafka.send(
            KafkaMessage{ topic: format!($topic), message: $message }
        ).await {
            log::error!("Failed to post to kafka: {:?}", e);
        }
    }
}

impl Kafka {
    pub async fn send(&self, message: KafkaMessage) -> anyhow::Result<()> {
        let mut tx = self.tx.clone();
        match tx.send(message).await {
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

fn start(mut rx: Receiver<KafkaMessage>) {
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
                        FutureRecord::to(message.topic.as_str())
                            .payload(message.message.as_str())
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