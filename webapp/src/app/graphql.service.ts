import { Injectable } from '@angular/core';
import { Apollo } from 'apollo-angular';
import gql from 'graphql-tag';
import { Observable } from 'apollo-client/util/Observable';

@Injectable({
  providedIn: 'root'
})
export class GraphqlService {

  constructor(
    private apollo: Apollo
  ) { }

  queryQuestions() {
    const QueryQuestions = gql`
      query {
        action(userID: 1) {
          questions {
            id
            title
            content
          }
        }
      }
    `;
    let obs = this.apollo.query({
      query: QueryQuestions
    });
    return obs
  }

  queryUsers() {
    const queryUsers = gql`
      query GetAllUserQuery{
        action(userID:2){
          users {
            id
            name
          }
        }
      }
    `;

    return this.apollo.query<any>({
      query: queryUsers
    });
  }
}

export class Question {
  id: Number
  title: string
  content: string
}
