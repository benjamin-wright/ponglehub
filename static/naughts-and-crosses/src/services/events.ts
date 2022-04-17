import { CloudEvent } from "cloudevents";

export type SendEvents = {
    send: (type: string, data: any) => void
    stop: () => void
}

type EventData = {
    type: string,
    data: any,
}

export function listen(
    receive: (type: string, data: any) => void
): SendEvents {
    let socket: WebSocket;
    let messages: EventData[] = [];
    let attempts = 0;
    let timeout = true;
    let connected = false;

    let send = (type: string, data: any) => {
        if (!connected) {
            messages.push({type, data});
            return;
        }

        const event = new CloudEvent({
            type,
            source: "web-cli",
            data
        });

        socket.send(event.toString());
    }

    let stop = () => {
        if (socket) {
            connected = false;
            socket.onclose = () => {};
            socket.onerror = () => {};
            socket.close();
            socket = undefined;
        }
    }

    let connect = () => {
        if (attempts > 3 && timeout) {
            timeout = false;
            setTimeout(connect, 5000);
            return;
        }

        attempts++;
        timeout = true
        
        stop();

        socket = new WebSocket("ws://ponglehub.co.uk/events")

        socket.onopen = () => {
            connected = true;
            attempts = 0;
            messages.forEach(m => send(m.type, m.data))
            messages = [];
        }

        socket.onmessage = (event: MessageEvent) => {
            const data = JSON.parse(event.data);
            receive(data.type, data.data);
        }
    
        socket.onclose = (event: Event) => {
            console.warn("socket connection closed");
            socket.onclose = () => {};
            socket.onerror = () => {};
            socket = undefined;
            connected = false;
            connect();
        };
    
        socket.onerror = (error: Event) => {
            console.warn(`Socket error: ${error}`);
            socket.onclose = () => {};
            socket.onerror = () => {};
            socket = undefined;
            connected = false;
            connect();
        }
    }

    connect();

    return {
        send: send,
        stop: stop
    };
}