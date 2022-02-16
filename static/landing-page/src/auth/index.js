import axios from "axios";

const landingPageUrl = "http://localhost:3000";
const authUrl = "http://localhost:4000";

export async function logOut() {
  const response = await axios.post(
    `${authUrl}/auth/logout`,
    {},
    { withCredentials: true }
  );
  return response.status == 204;
}

export async function getUserData() {
  const response = await axios.get(`${authUrl}/auth/user`, {
    withCredentials: true,
  });
  return response.data;
}

export function redirectToLogin() {
  window.location = `${authUrl}/auth/login?redirect=${window.location.toString()}`
}
