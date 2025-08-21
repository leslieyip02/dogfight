import p5 from "p5";

export interface Entity {
    position: p5.Vector;

    update: (x: number, y: number, theta: number) => void;
    draw: (instance: p5) => void;
};
