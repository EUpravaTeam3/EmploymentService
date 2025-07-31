import { ComponentFixture, TestBed } from '@angular/core/testing';

import { JobadComponent } from './jobad.component';

describe('JobadComponent', () => {
  let component: JobadComponent;
  let fixture: ComponentFixture<JobadComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ JobadComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(JobadComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
