use rocket::{ Data, Request, Rocket, request::{ FromRequest, Outcome } };
use rocket::fairing::{ Info, Fairing, Kind };
use rdkafka::producer::FutureProducer;

pub struct KafkaFairing {
    kafka: FutureProducer
}

pub struct Kafka {
    kafka: FutureProducer
}

#[derive(Debug)]
enum KafkaError {
    DodgyConnection,
}

impl<'a, 'r> FromRequest<'a, 'r> for Kafka {
    type Error = KafkaError;

    fn from_request(request: &'a Request<'r>) -> Outcome<Self, Self::Error> {
        Outcome::Success(Kafka{
            kafka: rdkafka::ClientConfig::new().create().expect("to work")
        })
    }
}

impl Fairing for KafkaFairing {
    fn info(&self) -> Info {
        Info {
            name: "KAFKA Connector",
            kind: Kind::Request | Kind::Launch
        }
    }

    fn on_attach(&self, rocket: Rocket) -> Result<Rocket, Rocket> {

    }

    fn on_request(&self, request: &mut Request, data: &Data) {
        request.
    }
}