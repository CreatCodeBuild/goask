import { Component, OnInit } from '@angular/core';
import { GraphqlService, Question } from '../graphql.service'

@Component({
  selector: 'app-questions',
  templateUrl: './questions.component.html',
  styleUrls: ['./questions.component.css']
})
export class QuestionsComponent implements OnInit {

  private questions: Array<any>

  constructor(
    private graphqlService: GraphqlService
  ) { 
    this.questions = new Array<Question>();
  }

  ngOnInit() {
    console.log("QuestionsComponent on init")
   
    console.log(this.graphqlService)

    this.questions = this.graphqlService.queryQuestions()
  }

}
