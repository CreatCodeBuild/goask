import { Injectable } from '@angular/core';
import { Apollo } from 'apollo-angular';
import gql from 'graphql-tag';

@Injectable({
  providedIn: 'root'
})
export class UserService {

  constructor(
    private apollo: Apollo
  ) { }

  getAllUsers() {
    return this.apollo.query<any>({
      query: GetAllUserQuery
    });
  }
}

const GetAllUserQuery = gql`
query GetAllUserQuery{
  action(userID:2){
    users {
      id
      name
    }
  }
}
`;
