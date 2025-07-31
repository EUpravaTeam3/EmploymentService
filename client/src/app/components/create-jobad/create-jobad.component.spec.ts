import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateJobadComponent } from './create-jobad.component';

describe('CreateJobadComponent', () => {
  let component: CreateJobadComponent;
  let fixture: ComponentFixture<CreateJobadComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ CreateJobadComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreateJobadComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
