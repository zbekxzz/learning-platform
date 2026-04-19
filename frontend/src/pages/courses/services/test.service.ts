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

  submitTest(payload: TestSubmitPayload): Observable<TestSubmitResponse> {
    return this.api.post<TestSubmitResponse>(`/back/tests/submit`, payload);
  }
}
