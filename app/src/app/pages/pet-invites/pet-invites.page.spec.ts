import { ComponentFixture, TestBed } from '@angular/core/testing';
import { IonicModule } from '@ionic/angular';

import { PetInvitesPage } from './pet-invites.page';

describe('PetInvitesPage', () => {
  let component: PetInvitesPage;
  let fixture: ComponentFixture<PetInvitesPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PetInvitesPage, IonicModule],
    }).compileComponents();

    fixture = TestBed.createComponent(PetInvitesPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
