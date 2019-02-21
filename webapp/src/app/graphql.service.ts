import { Injectable } from '@angular/core';
import { Apollo } from 'apollo-angular';
import gql from 'graphql-tag';

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
        action(userID: "1") {
          questions {
            id
            title
            content
          }
        }
      }
    `;
    let obs = this.apollo.query<Data>({
      query: QueryQuestions
    });
    return obs
  }

  queryQuestionDetail(questionID: string) {
    const QueryQuestions = gql`
      query ($questionID: ID!) {
        action(userID: "1") {
          question(id: $questionID) {
            id
            title
            content
            author {
              id
              name
            }
            answers {
              id
              content
              author {
                id
                name
              }
            }
          }
        }
      }
    `;
    return this.apollo.query<Data>({
      query: QueryQuestions,
      variables: {
        "questionID": questionID,
      }
    });
  }

  queryUsers() {
    const queryUsers = gql`
      query GetAllUserQuery{
        action(userID:"1"){
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

export interface Data {
  action: Action
}

interface Action {
  questions?: Question[]
  question?: Question
}

export interface  Question {
  id: string
  title: string
  content: string
  author: Author
  answers: Answer[]
}

interface  Author {
  name: string
}

export interface Answer {
  id: string
  content: string
  author: Author
}
