import { jwtDecode } from "jwt-decode";
import { $authHost, $host } from ".";

export const registration = async (email: string, login: string, password: string) => {
  const {data} = await $host.post('signup', {email, login, password})
  return jwtDecode(data.access_token)
}

export const auth = async (login_or_email: string, password: string) => {
  const {data} = await $host.post('signin', {login_or_email, password})
  localStorage.setItem('token', data.access_token)
  return jwtDecode(data.access_token)
}

export const check = async () => {
  const response = await $host.post('signup', {})
  return response
}

