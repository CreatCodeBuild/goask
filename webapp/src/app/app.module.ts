import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule }   from '@angular/forms';
import { ApolloModule, APOLLO_OPTIONS } from 'apollo-angular';
import { HttpLinkModule, HttpLink } from 'apollo-angular-link-http';
import { InMemoryCache } from 'apollo-cache-inmemory';

import { AppComponent } from './app.component';
import { environment } from 'src/environments/environment';
import { UsersComponent } from './users/users.component';
import { Routes, RouterModule } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';
import { QuestionsComponent } from './questions/questions.component';
import { QuestionDetailComponent } from './question-detail/question-detail.component';
import { NavComponent } from './nav/nav.component';
import { AnswerDetailComponent } from './answer-detail/answer-detail.component';
import { UserDetailComponent } from './user-detail/user-detail.component';
import { UserSummaryComponent } from './user-summary/user-summary.component';
import { SigninComponent } from './signin/signin.component';
import { SignupComponent } from './signup/signup.component';
import { UserMeComponent } from './user-me/user-me.component';

const routes: Routes = [
  {
    path: 'users',
    component: UsersComponent
  },
  {
    path: 'question/:id',
    component: QuestionDetailComponent,
  },
  {
    path: 'questions',
    component: QuestionsComponent,
  },
  {
    path: 'me',
    component: UserMeComponent,
  },
  {
    path: 'signin',
    component: SigninComponent,
  },
  {
    path: 'signup',
    component: SignupComponent,
  }
];

@NgModule({
  declarations: [
    AppComponent,
    UsersComponent,
    QuestionsComponent,
    QuestionDetailComponent,
    NavComponent,
    AnswerDetailComponent,
    UserDetailComponent,
    UserSummaryComponent,
    SigninComponent,
    SignupComponent,
    UserMeComponent
  ],
  imports: [
    BrowserModule,
    ApolloModule,
    HttpClientModule,
    HttpLinkModule,
    RouterModule.forRoot(routes),
    FormsModule,
  ],
  providers: [
    {
      provide: APOLLO_OPTIONS,
      useFactory(httpLink: HttpLink) {
        return {
          cache: new InMemoryCache(),
          link: httpLink.create({
            uri: environment.graphQLEndpoint
          })
        };
      },
      deps: [HttpLink]
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
