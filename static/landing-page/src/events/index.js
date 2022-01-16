import axios from "axios";

export default class {
  constructor(url) {
    this.url = url;
  }

  async logOut() {
    const response = await axios.post(
      `${this.url}/auth/logout`,
      {},
      { withCredentials: true }
    );
    return response.status == 204;
  }

  async getUserData() {
    const response = await axios.get(`${this.url}/auth/user`, {
      withCredentials: true,
    });
    return response.data;
  }
}
