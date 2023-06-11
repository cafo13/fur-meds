import {
  HttpClient,
  HttpErrorResponse,
  HttpHeaders,
} from '@angular/common/http';
import { Injectable } from '@angular/core';

import { Observable, catchError, throwError } from 'rxjs';

import { environment } from 'src/environments/environment';
import { Pet, PetShareInvite } from '../types/types';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  apiBaseDomain = environment.apiEndpoint;

  constructor(private http: HttpClient) {}

  private getAccessToken(): string | undefined {
    const currentUser = JSON.parse(localStorage.getItem('user')!);
    return currentUser?.stsTokenManager?.accessToken ?? undefined;
  }

  private getHeaders(): HttpHeaders {
    return new HttpHeaders({
      Authorization: `Bearer ${this.getAccessToken()}`,
    });
  }

  private handleError(error: HttpErrorResponse) {
    if (error.status === 0) {
      console.error('An error occurred:', error.error);
    } else {
      console.error(
        `Api returned code ${error.status}, body was: `,
        error.error
      );
    }

    return throwError(
      () => new Error(`Error, please try again later: ${error.error}`)
    );
  }

  // get the full pet list
  public getPets(): Observable<Pet[]> {
    const headers = this.getHeaders();
    return this.http
      .get<Pet[]>(`${this.apiBaseDomain}/pets`, { headers })
      .pipe(catchError(this.handleError));
  }

  // add a pet, returns the new full pet list
  public addPet(pet: Partial<Pet>): Observable<Pet[]> {
    const headers = this.getHeaders();
    return this.http
      .post<Pet[]>(`${this.apiBaseDomain}/pet`, pet, { headers })
      .pipe(catchError(this.handleError));
  }

  // update a pet, returns the new full pet list
  public updatePet(pet: Partial<Pet>): Observable<Pet[]> {
    const headers = this.getHeaders();
    return this.http
      .put<Pet[]>(`${this.apiBaseDomain}/pet`, pet, { headers })
      .pipe(catchError(this.handleError));
  }

  // delete a pet, returns the new full pet list
  public deletePet(uuid: string): Observable<Pet[]> {
    const headers = this.getHeaders();
    return this.http
      .delete<Pet[]>(`${this.apiBaseDomain}/pet/${uuid}`, { headers })
      .pipe(catchError(this.handleError));
  }

  // invite another user to one of your pets
  public inviteUserToSharedPet(petUuid: string, userMailToInvite: string) {
    const headers = this.getHeaders();
    return this.http
      .post(
        `${this.apiBaseDomain}/pet/share/invite`,
        { petUuid, userMailToInvite },
        { headers }
      )
      .pipe(catchError(this.handleError));
  }

  // get the full list of open pet share invites for user
  public getPetShareInvites(): Observable<PetShareInvite[]> {
    const headers = this.getHeaders();
    return this.http
      .get<PetShareInvite[]>(`${this.apiBaseDomain}/pet/share/invites`, {
        headers,
      })
      .pipe(catchError(this.handleError));
  }

  // accept an invite to a shared pet
  public acceptInviteToSharedPet(petUuid: string) {
    const headers = this.getHeaders();
    return this.http
      .post(`${this.apiBaseDomain}/pet/share/accept`, { petUuid }, { headers })
      .pipe(catchError(this.handleError));
  }

  // deny an invite to a shared pet
  public denyInviteToSharedPet(petUuid: string) {
    const headers = this.getHeaders();
    return this.http
      .post(`${this.apiBaseDomain}/pet/share/deny`, { petUuid }, { headers })
      .pipe(catchError(this.handleError));
  }
}
