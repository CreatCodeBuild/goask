import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
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
  }
];

@NgModule({
  declarations: [
    AppComponent,
    UsersComponent,
    QuestionsComponent,
    QuestionDetailComponent,
    NavComponent
  ],
  imports: [
    BrowserModule,
    ApolloModule,
    HttpClientModule,
    HttpLinkModule,
    RouterModule.forRoot(routes)
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
