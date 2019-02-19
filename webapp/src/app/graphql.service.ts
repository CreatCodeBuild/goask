import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class GraphqlService {

  constructor() { }

  queryQuestions(): Array<Question> {
    return [{
      id: 1,
      title: "title 1",
      content: "content 1"
    },{
      id: 2,
      title: "title 2",
      content: "content 2"
    }]
  }
}

export class Question {
  id: Number
  title: string
  content: string
}
