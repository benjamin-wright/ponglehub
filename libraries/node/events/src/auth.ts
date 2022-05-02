import axios from "axios";
import { UserData } from './storage';

export class Auth {
  private host: string;

  constructor(host: string) {
    this.host = host;
  }

  async load(): Promise<UserData> {
    try {

      const response = await axios.get(`http://${this.host}/auth/user`, {
        withCredentials: true
      });
      
      if (typeof response.data.name !== "string") {
        throw new Error(`parsing error, name property not found on "${JSON.stringify(response.data)}"`);
      }
  
      const data: UserData = {
        name: response.data.name
      };
  
      return data;
    } catch (error) {
      if (error.response.status == 401) {
        return null;
      }

      throw error;
    }
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
