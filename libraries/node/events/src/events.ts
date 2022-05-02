import { CloudEvent } from "cloudevents";

export class Events {
    private host: string;
    private socket: WebSocket;

    constructor(host: string) {
        this.host = host;
        this.socket = null;
    }

    start(receive: (type: string, data: any) => void, closed: () => void): Promise<void> {
        const socket = new WebSocket(`ws://${this.host}/events`);

        return new Promise((resolve, reject) => {
            socket.onopen = () => {
                this.socket = socket;
                resolve();
            }
    
            socket.onmessage = (event: MessageEvent) => {
                const data = JSON.parse(event.data);
                receive(data.type, data.data);
            }
        
            socket.onclose = () => {
                console.warn("socket connection closed");
                socket.onclose = () => {};
                socket.onerror = () => {};
                this.socket = null;
                closed();
            };
        
            socket.onerror = (error: Event) => {
                console.warn('Socket error', error);
                socket.onclose = () => {};
                socket.onerror = () => {};
                this.socket = null;
                reject(error);
            }
        });
    }

    send(type: string, data: any): void {
        if (!this.socket) {
            throw new Error(`Tried to send message ${type} to closed websocket`);
        }

        const event = new CloudEvent({
            type,
            source: "web-cli",
            data
        });

        this.socket.send(event.toString());
    }

    stop(): void {
        if (this.socket) {
            this.socket.onclose = () => {};
            this.socket.onerror = () => {};
            this.socket.close();
            this.socket = null;
        }
    }
}
