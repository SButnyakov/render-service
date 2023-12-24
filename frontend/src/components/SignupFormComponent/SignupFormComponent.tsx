import React, { FormEvent, useState } from "react"
import { registration } from "../../http/AuthAPI"
import { observer } from "mobx-react-lite"
import { useNavigate } from "react-router-dom"
import { AxiosError } from "axios"

import styles from './SignupFormComponent.module.css'
import { useForm } from "react-hook-form"

type FormData = {
  login: string;
  email: string;
  password: string;
  passwordRepeat: string;
};

/* TODO: Доделать валидацию*/

export const SignupForm = observer(() => {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<FormData>();

  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')
  const [email, setEmail] = useState('')

  const route = useNavigate()

  const onSubmit = (_: any) => {
    signUp()
  }

  const signUp = async () => {
    await registration(email, login, password)
      .then(_ => {
        route('/signin')
      })
      .catch((error: AxiosError) => {

      })
  }

  return(
    <div >
      <form onSubmit={handleSubmit(onSubmit)}>
      <div>
          <label htmlFor="login">Login:</label>
          <input 
            id="login"
            {...register('login', {
              required: 'Enter login',
              pattern: {
                value: /^[a-zA-Z\d]{4,15}$/,
                message: 'Логин должен состоять из 4-15 символов (латинские буквы и цифры)',
              }
            })}
          />
        </div>

        <div>
          <label htmlFor="email">Email:</label>
          <input 
            name="email" 
            type="text" 
            value={email} 
            onChange={e => {setEmail(e.target.value)}}
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

        <div className={styles.errorBlockMessage}>
          {errors.login && (
            <p>{errors.login.message}</p>
          )}

        </div>
        

        <button>Register</button>
        <button onClick={() => {route('/signin')}}>Log In</button>
      </form>
    </div>
  )
})

export default SignupForm