import axios from "axios"

// Для обычных запросов, которые не требуют авторизации
const $host = axios.create({
  baseURL: process.env.REACT_APP_API_AUTH_URL
})

// Тут для каждого запроса будет еще подставляться JWT токен
const $authHost = axios.create({
  baseURL: process.env.REACT_APP_API_SERVER_URL
})

const authInterceptor = (config: any) => {
  config.headers.authorization = `Bearer ${localStorage.getItem('token')}`
  return config
}

const invalidTokenInterceptor = async (error: any) => {
  const originalReq = {...error.config}
  originalReq._isRetry = true

  if (error.response?.status === 401) {
    try {
      const currentRefreshToken = localStorage.getItem('refresh_token')
      const {data} = await $host.put('refresh', {refresh_token: currentRefreshToken})

      localStorage.setItem('token', data.access_token)
      localStorage.setItem('refresh_token', data.refresh_token)

      return $authHost.request(originalReq)
    }
    catch (error) {
      console.error('[AUTH_ERROR]: Error with refresh token!')
    }
  }

  throw error
}

const validTokenInterceptor = (config: any) => config

$authHost.interceptors.request.use(authInterceptor)

$authHost.interceptors.response.use(validTokenInterceptor, invalidTokenInterceptor)

export {
  $host,
  $authHost
}

