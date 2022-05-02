import axios from "axios";

export class Auth {
  private host: string;

  constructor(host: string) {
    this.host = host;
  }

  async logOut(): Promise<any> {
    const response = await axios.post(
      `http://${this.host}/auth/logout`,
      {},
      { withCredentials: true }
    );

    if (response.status != 204) {
      throw new Error(`failed signing out: status code ${response.status}`);
    };
  }

  logIn() {
    window.location.href = `http://${this.host}/auth/login?redirect=${window.location.toString()}`;
  }
}
