import { Component, OnInit } from '@angular/core';
import { GraphqlService, Question } from '../graphql.service'
import {map} from 'rxjs/operators';

@Component({
  selector: 'app-questions',
  templateUrl: './questions.component.html',
  styleUrls: ['./questions.component.css']
})
export class QuestionsComponent implements OnInit {

  private questions: Array<Question>

  constructor(
    private graphqlService: GraphqlService
  ) { 
    this.questions = new Array<Question>();
  }

  ngOnInit() {
    let self = this
    this.graphqlService.queryQuestions().subscribe({
      next(x) { 
        self.questions = self.questions.concat(x.data.action.questions)
      }
    })
  }

}
