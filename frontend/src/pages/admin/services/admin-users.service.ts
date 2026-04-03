import { Injectable } from '@angular/core';
import { ApiService } from '../../../shared/services/api.service';
import { Observable } from 'rxjs';
import { UserProfile } from '../../profile/models/profile.model';
import { UserPayload } from '../models/admin-user.model';

@Injectable({
  providedIn: 'root'
})
export class AdminUsersService {

  constructor(private api: ApiService) {}

  getUsers(): Observable<UserProfile[] | any> {
    return this.api.get('/back/users/');
  }

  createUser(data: UserPayload): Observable<UserProfile> {
    return this.api.post('/back/users/', data);
  }

  updateUser(id: number, data: UserPayload): Observable<UserProfile> {
    return this.api.put(`/back/users/${id}`, data);
  }

  deleteUser(id: number): Observable<any> {
    return this.api.delete(`/back/users/${id}`);
  }
}
