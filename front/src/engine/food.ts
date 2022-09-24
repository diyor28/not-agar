import Entity from "./entity";
import {FoodData} from "../client";

export default class Food extends Entity {
	id: number

	constructor(data: FoodData) {
		super(data);
		this.id = data.id;
	}
}
