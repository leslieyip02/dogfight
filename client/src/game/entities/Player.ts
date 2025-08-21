import p5 from "p5";
import type { Entity } from "./Entity";

class Player implements Entity {
  position: p5.Vector;
  theta: number;
  username: string;
  roll: number;

  constructor(username: string, x: number, y: number, theta: number) {
    this.username = username;
    this.position = new p5.Vector(x, y);
    this.theta = theta;
    this.roll = 0;
  }

  update = (x: number, y: number, theta: number) => {
    this.position.x = x;
    this.position.y = y;

    this.roll = Math.sign(theta - this.theta);
    this.theta = theta;
  };

  draw = (instance: p5) => {
    instance.push();
    
    instance.translate(
      -this.position.x + window.innerWidth / 2,
      -this.position.y + window.innerHeight / 2
    );

    instance.translate(this.position.x, this.position.y);
  
    instance.push();
    instance.rotate(this.theta);
    instance.fill("#ffffff");
    instance.triangle(
      -40, 0,
      40, 40,
      40, -40,
    );
    instance.pop();

    instance.noFill();
    instance.stroke("#ffffff");
    instance.circle(0, 0, 80);
  
    instance.strokeWeight(1);
    instance.textAlign(instance.CENTER);
    instance.rectMode(instance.CENTER);
    instance.text(this.username, 0, -65);
    instance.pop();
  };
};

export default Player;
