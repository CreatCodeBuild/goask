import { Component, OnInit } from '@angular/core';
import { GraphqlService, Question } from '../graphql.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-question-detail',
  templateUrl: './question-detail.component.html',
  styleUrls: ['./question-detail.component.css']
})
export class QuestionDetailComponent implements OnInit {

  private question: Question;
  private questionID: string;

  constructor(
    private graphqlService: GraphqlService,
    private route: ActivatedRoute // todo? how to use it?
  ) { 

  }

  ngOnInit() {
    this.route.paramMap.subscribe(async (paramMap) => {
      this.questionID = paramMap.get("id")
      let result = await this.graphqlService.queryQuestionDetail(this.questionID).toPromise()
      this.question = result.data.action.question
    })  
  }

}
