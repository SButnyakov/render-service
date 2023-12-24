import { $host } from "."
import { AxiosError } from "axios"

export const registration = async (email: string, login: string, password: string) => {
  await $host.post('signup', {"email": email, "login": login, "password": password})
    .catch(err => {
      console.error(err)
    })
}

export const auth = async (login_or_email: string, password: string) => {
  await $host.post('signin', {login_or_email, password})
      .then(({data}) => {
        localStorage.setItem('token', data.access_token)
        localStorage.setItem('refresh_token', data.refresh_token)
      })
      .catch((error: AxiosError) => {
        throw error
      })
}
