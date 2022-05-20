export class Events {
    private host: string;
    private socket: WebSocket | null;

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
                const parsed = JSON.parse(event.data);

                console.log("parsed data", parsed.data);
                const data = typeof(parsed.data) === "string" ? JSON.parse(parsed.data) : parsed.data;
                console.log("double-parsed", data);
                
                receive(parsed.type, data);
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

        this.socket.send(JSON.stringify({type,data}));
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
