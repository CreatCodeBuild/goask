import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private user: User

  constructor() {
    this.user = {
      id: "", // empty string represents an unlogged in user
      name: "游客"
    }
  }

  set User(user: User) {
    this.user = user
  }

  current(): User {
    return this.user
  }
}

export interface User {
  id: string
  name: string
}
