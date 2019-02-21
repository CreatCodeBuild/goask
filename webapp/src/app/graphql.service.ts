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

  queryQuestions(userID: string) {
    const QueryQuestions = gql`
      query ($userID: ID!) {
        action(userID: $userID) {
          questions {
            id
            title
            content
          }
        }
      }
    `;
    let obs = this.apollo.query<Data>({
      query: QueryQuestions,
      variables: {
        "userID": userID,
      }
    });
    return obs
  }

  queryQuestionDetail(userID: string, questionID: string) {
    const QueryQuestions = gql`
      query ($userID: ID!, $questionID: ID!) {
        action(userID: $userID) {
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
        "userID": userID,
      }
    });
  }

  queryUser(actionTakerID: string, userID: string) {
    const queryUsers = gql`
      query ($actionTakerID: ID!, $userID: ID!) {
        action(userID: $actionTakerID) {
          user(id: $userID) {
            id
            name
            questions {
              id
              answers {
                id
              }
            }
          }
        }
      }
    `;

    return this.apollo.query<any>({
      query: queryUsers,
      variables: {
        "actionTakerID": actionTakerID,
        "userID": userID,
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
