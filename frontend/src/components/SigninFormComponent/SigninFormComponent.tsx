import React, { FormEvent, useState } from "react"
import { auth } from "../../http/AuthAPI"
import { observer } from "mobx-react-lite"

import { useStore } from "../../hooks/useStore"
import { useNavigate } from "react-router-dom"

import styles from './SigninFormComponent.module.css'
import { AxiosError } from "axios"
import { SigninResponseCodes } from "../../http/httpTypes"

export const SigninForm = observer(() => {
  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')

  const [errorMessage, setErrorMessage] = useState('')

  const store = useStore()

  const route = useNavigate()

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    signIn()
  }

  const signIn = async () => {
    await auth(login, password)
      .then(_ => {
        store.userStore.setIsAuth(true)

        route('/upload')
      })
      .catch((error: AxiosError) => {
        const {response} = error

        if (response?.status === SigninResponseCodes.INTERNAL_SERVER_ERROR) {
          setErrorMessage('Не удалось проверить данные для входа')
        }

        if (response?.status === SigninResponseCodes.INVALID_CREDENTIALS) {
          setErrorMessage('Неверный логин или пароль')

          setLogin('')
          setPassword('')
        }
      })
  }

  return(
    <div >
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="login">Login:</label>
          <input 
            name="login" 
            type="text" 
            value={login} 
            onChange={e => {setLogin(e.target.value)}}
          />
        </div>

        <div>
          <label htmlFor="password">Password:</label>
          <input 
            name="password" 
            type="password" 
            value={password} 
            onChange={e => {setPassword(e.target.value)}}
          />
        </div>
        <button disabled={!login || !password}>Log In</button>
        <button onClick={() => {route('/signup')}}>Register</button>
      </form>
      <div className={styles.errorBlockMessage}>
        {errorMessage}
      </div>
    </div>
  )
})

export default SigninForm