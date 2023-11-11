import axios from "axios";

// Для обычных запросов, которые не требуют авторизации
const $host = axios.create({
  baseURL: process.env.REACT_APP_API_URL
})

// Тут для каждого запроса будет еще подставляться JWT токен
const $authHost = axios.create({
  baseURL: process.env.REACT_APP_API_URL
})

const authInterceptor = (config: any) => {
  config.headers.authorization = `Bearer ${localStorage.getItem('token')}`
  return config
}

$authHost.interceptors.request.use(authInterceptor)

export {
  $host,
  $authHost
}

