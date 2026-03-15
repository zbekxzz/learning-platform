import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { LoginResponse } from '../../core/auth/auth.model';
import { UserModel } from '../models/user.model';
import { Nullable } from 'primeng/ts-helpers';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  _userData$: BehaviorSubject<Nullable<LoginResponse | UserModel>> = new BehaviorSubject<
    Nullable<LoginResponse | UserModel>
  >(null);
  userData: Nullable<LoginResponse | UserModel> = this._userData$.getValue();
  userData$: Observable<Nullable<LoginResponse | UserModel>> = this._userData$.asObservable();

  setUserData(data: LoginResponse | UserModel) {
    this._userData$.next(data);
  }
}
