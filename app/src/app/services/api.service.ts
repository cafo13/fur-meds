import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';

import { Observable, catchError, throwError } from 'rxjs';

import { environment } from 'src/environments/environment';
import { Pet } from '../types/types';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  apiBaseDomain = environment.apiEndpoint;

  constructor(private http: HttpClient) {}

  private handleError(error: HttpErrorResponse) {
    if (error.status === 0) {
      console.error('An error occurred:', error.error);
    } else {
      console.error(
        `Backend returned code ${error.status}, body was: `,
        error.error
      );
    }

    return throwError(
      () => new Error('Something bad happened; please try again later.')
    );
  }

  // get the full pet list
  public getPetList(): Observable<Pet[]> {
    return this.http
      .get<Pet[]>(`${this.apiBaseDomain}/pets`)
      .pipe(catchError(this.handleError));
  }

  // add a pet, returns the new full pet list
  public addPet(pet: Partial<Pet>): Observable<Pet[]> {
    return this.http
      .post<Pet[]>(`${this.apiBaseDomain}/pet`, pet)
      .pipe(catchError(this.handleError));
  }

  // update a pet
  public updatePet(pet: Partial<Pet>): Observable<Pet> {
    return this.http
      .put<Pet>(`${this.apiBaseDomain}/pet`, pet)
      .pipe(catchError(this.handleError));
  }

  // delete a pet
  public deletePet(uuid: string): Observable<unknown> {
    return this.http
      .delete(`${this.apiBaseDomain}/pet/${uuid}`)
      .pipe(catchError(this.handleError));
  }
}
