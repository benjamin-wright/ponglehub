use std::{time::Duration, time::SystemTime};

use rocket::{Data, Request, Response, fairing::{Fairing, Info, Kind}};

pub struct Logger {}
#[derive(Copy, Clone)]
struct TimerStart(Option<SystemTime>);

impl Fairing for Logger {
    fn info(&self) -> Info {
        Info{
            name: "Logging Middleware",
            kind: Kind::Request | Kind::Response
        }
    }

    fn on_request(&self, request: &mut Request, _data: &Data) {
        request.local_cache(|| TimerStart(Some(SystemTime::now())));

        let method = request.method().as_str();
        let path = request.uri().path();

        println!("HTTP {} -> {}", method, path);
    }

    fn on_response(&self, request: &Request, response: &mut Response) {
        let start_time = request.local_cache(|| TimerStart(None));
        let mut dur: f32 = 0.0;
        if let Some(Ok(duration)) = start_time.0.map(|st| st.elapsed()) {
            if duration.as_millis() > 5 {
                dur = (duration.as_secs() * 1000 + duration.subsec_millis() as u64) as f32;
            } else {
                dur = (duration.as_secs() * 1000) as f32 + duration.subsec_micros() as f32 / 1000.0;
            }
        }

        let method = request.method().as_str();
        let path = request.uri().path();
        let status = response.status().code;

        println!("HTTP {} <- {} [{}] ({}ms)", method, path, status, dur);
    }
}