import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { CV } from 'src/app/model/cv';
import { Diploma } from 'src/app/model/diploma';

@Component({
  selector: 'app-cv',
  templateUrl: './cv.component.html',
  styleUrls: ['./cv.component.css']
})
export class CvComponent implements OnInit {

  cv: CV = {
    citizen_ucn: '',
    description: '',
    name: '',
    email: '',
    work_experience: [],
    education: []
  }

  constructor(private http: HttpClient) {}

  ngOnInit() {

    var ucn = localStorage.getItem("eupravaUcn")
    var name = localStorage.getItem("eupravaName") + " " + localStorage.getItem("eupravaSurname")
    var email = localStorage.getItem("eupravaEmail")
    if (ucn && name && email){
    this.cv.citizen_ucn = ucn
    this.cv.name = name
    this.cv.email = email
    }
      
    this.http.get<CV>('http://localhost:8000/resume/citizen/' + ucn, { withCredentials: true })
      .subscribe(cv => {
      this.cv = cv
      }, err => {
        console.log(err);
        alert(err)
      });

  }

  onEditResume(){
    if (this.cv.education.length < 1){
      alert("you must generate your diploma!")
      return
    }

    this.http.post('http://localhost:8000/resume/' + this.cv.citizen_ucn, this.cv, { withCredentials: true})
    .subscribe(res => console.log(res), err => alert(err))
  }

  onGenerateDiploma(){

    var diplomas: Diploma[] = [
      { _id: "", institution_id: "", institution_name:"Faculty of technical sciences", institution_type:"college",
        average_grade: 7, ucn: this.cv.citizen_ucn
      },
      { _id: "", institution_id: "", institution_name:"Faculty of matehmatics", institution_type:"college",
        average_grade: 9, ucn: this.cv.citizen_ucn
      },
      { _id: "", institution_id: "", institution_name:"Faculty of agriculture", institution_type:"college",
        average_grade: 8, ucn: this.cv.citizen_ucn
      },
      { _id: "", institution_id: "", institution_name:"Jovan Jovanovic Zmaj", institution_type:"high school",
        average_grade: 5, ucn: this.cv.citizen_ucn
      }
    ]

    this.cv.education.push(diplomas[this.randomIntFromInterval(0, 3)])
  }

  randomIntFromInterval(min: number, max: number) {
  return Math.floor(Math.random() * (max - min + 1) + min);
}

}
