import { Injectable } from '@angular/core';
import { ApiService } from '../../../shared/services/api.service';
import { Observable } from 'rxjs';
import { TestResponse, TestSubmitPayload, TestSubmitResponse } from '../models/test.model';

@Injectable({
  providedIn: 'root'
})
export class TestService {

  constructor(private api: ApiService) {}

  startTest(moduleId: number): Observable<TestResponse> {
    return this.api.get<TestResponse>(`/back/tests/module/${moduleId}/start`);
  }

  startChapterTest(chapterId: number): Observable<TestResponse> {
    return this.api.get<TestResponse>(`/back/tests/chapter/${chapterId}/start`);
  }

  startFinalTest(courseId: number): Observable<TestResponse> {
    return this.api.get<TestResponse>(`/back/tests/course/${courseId}/final/start`);
  }

  submitTest(payload: TestSubmitPayload): Observable<TestSubmitResponse> {
    return this.api.post<TestSubmitResponse>(`/back/tests/submit`, payload);
  }

  createFullTest(payload: any): Observable<any> {
    return this.api.post<any>(`/back/tests/create-full`, payload);
  }

  getTestDetails(testId: number): Observable<any> {
    return this.api.get<any>(`/back/tests/${testId}/details`);
  }

  deleteTest(testId: number): Observable<any> {
    return this.api.delete<any>(`/back/tests/${testId}`);
  }
}

