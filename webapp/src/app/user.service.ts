import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class UserService {

  constructor() { }

  current(): User {
    return {
      id: "1", // empty string represents an unlogged in user
      name: "游客"
    }
  }
}

export interface User {
  id: string
  name: string
}
