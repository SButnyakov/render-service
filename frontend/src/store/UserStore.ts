import { makeAutoObservable } from "mobx";

export class UserStore {
  private _isAuth: boolean;
  private _user;

  constructor() {
    this._isAuth = false;
    this._user = {};

    makeAutoObservable(this);
  }

  public setIsAuth(isAuth: boolean) {
    this._isAuth = isAuth;
  }

  public setUser(user: any) {
    this._user = user
  }

  get isAuth() {
    return this._isAuth
  }

  get user() {
    return this._user
  }
}

export default new UserStore()