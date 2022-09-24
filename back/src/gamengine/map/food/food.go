package food

import (
	"github.com/diyor28/not-agar/src/gamengine/map/entity"
	"github.com/diyor28/not-agar/src/utils"
)

type EntityInterface interface {
	getX() float32
	getY() float32
}

type Food struct {
	*entity.Entity
	Color [3]uint8 `json:"color"`
}

type Foods struct {
	Food   []*Food
	lastId entity.Id
}

func (f *Foods) New(x float32, y float32, weight float32) *Food {
	f.lastId++
	id := f.lastId
	food := &Food{
		Entity: &entity.Entity{
			Id:     id,
			X:      x,
			Y:      y,
			Weight: weight,
		},
		Color: utils.RandomColor(),
	}
	f.Append(food)
	return food
}

func (f *Foods) Copy() []*Food {
	result := make([]*Food, len(f.Food))
	copy(result, f.Food)
	return result
}

func (f *Foods) Append(food *Food) {
	f.Food = append(f.Food, food)
}

func (f *Foods) RemoveFoodById(id entity.Id) {
	for i, food := range f.Food {
		if food.Id == id {
			f.Food = append(f.Food[:i], f.Food[i+1:]...)
		}
	}
}

func (f *Foods) RemoveFoodAt(index int) {
	f.Food = append(f.Food[:index], f.Food[index+1:]...)
}

func (f *Foods) Last() *Food {
	return f.Food[len(f.Food)-1]
}

func (f *Foods) Len() int {
	return len(f.Food)
}

func (f *Foods) Closest(player entity.Interface, kClosest int) []*Food {
	foodCopy := f.Copy()
	totalNumFood := len(foodCopy)
	foodDistances := make(map[entity.Id]float32, totalNumFood)
	for _, f := range foodCopy {
		foodDistances[f.Id] = utils.Distance(player.GetX(), f.X, player.GetY(), f.Y)
	}
	numResults := kClosest
	if kClosest > totalNumFood {
		numResults = totalNumFood
	}
	for i := 0; i < numResults; i++ {
		var minIdx = i
		for j := i + 1; j < totalNumFood; j++ {
			if foodDistances[foodCopy[j].Id] < foodDistances[foodCopy[minIdx].Id] {
				minIdx = j
			}
		}
		foodCopy[i], foodCopy[minIdx] = foodCopy[minIdx], foodCopy[i]
	}
	return foodCopy[:numResults]
}
