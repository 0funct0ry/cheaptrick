package data

import (
	"fmt"
	"math/rand"
	"time"
)

var promptTemplates = []string{
	"What is %s?",
	"Explain %s in simple terms.",
	"How does %s work?",
	"Give a short summary of %s.",
	"Why is %s important?",
	"Tell me about %s.",
	"What are the benefits of %s?",
	"What are the risks of %s?",
	"Give a beginner explanation of %s.",
	"How would you describe %s to a student?",
}

var devTemplates = []string{
	"How do I implement %s in Go?",
	"What is the purpose of %s in programming?",
	"Explain %s with an example.",
	"When should developers use %s?",
	"What problem does %s solve?",
}

var emailTemplates = []string{
	"Write a short email about %s.",
	"How should I professionally ask about %s in an email?",
	"Draft a quick email regarding %s.",
}

var topics = []string{
	"machine learning",
	"docker containers",
	"REST APIs",
	"golang channels",
	"blockchain technology",
	"quantum computing",
	"healthy eating",
	"table tennis training",
	"stress management",
	"travel planning",
	"cybersecurity",
	"cloud computing",
	"neural networks",
	"renewable energy",
	"electric vehicles",
	"nutrition",
	"sleep quality",
	"photography basics",
	"music theory",
	"open source software",
	"database indexing",
	"distributed systems",
	"cryptography",
	"linux system administration",
	"functional programming",
}

var responsePool = []string{
	"%s is an important concept that appears in many modern systems. Understanding its basic principles can help people make better technical or practical decisions.",
	"%s refers to a field or idea that has grown rapidly in recent years. It combines theory and practical techniques to solve real-world problems.",
	"In simple terms, %s involves methods or practices designed to improve efficiency, understanding, or performance in a given domain.",
	"Learning about %s usually starts with the core principles and then moves toward practical examples and applications.",
	"%s is widely used across industries today and continues to evolve as technology and research advance.",
}

var longResponses = []string{
	"%s is a broad topic with many applications. At a high level it focuses on solving practical problems through structured approaches. Many professionals begin by learning the core concepts and then gradually applying them to real-world scenarios.",
	"The concept of %s has developed over time as researchers and practitioners explored better ways to address complex challenges. Today it plays a role in technology, science, and everyday decision making.",
	"When people first study %s they usually encounter the fundamental theory before moving on to hands-on experimentation. This combination of theory and practice is what makes the field both challenging and rewarding.",
}

func randomTemplate() string {
	sets := [][]string{
		promptTemplates,
		devTemplates,
		emailTemplates,
	}

	set := sets[rand.Intn(len(sets))]
	return set[rand.Intn(len(set))]
}

func randomResponse(topic string) string {

	if rand.Intn(2) == 0 {
		t := responsePool[rand.Intn(len(responsePool))]
		return fmt.Sprintf(t, topic)
	}

	t := longResponses[rand.Intn(len(longResponses))]
	return fmt.Sprintf(t, topic)
}

func GenerateTextPromptDataset(n int) map[string]string {

	_ = rand.New(rand.NewSource(time.Now().UnixNano()))

	data := make(map[string]string)

	for len(data) < n {

		topic := topics[rand.Intn(len(topics))]
		template := randomTemplate()

		prompt := fmt.Sprintf(template, topic)
		response := randomResponse(topic)

		data[prompt] = response
	}

	return data
}

//func main() {
//
//	dataset := GenerateTextPromptDataset(2000)
//
//	for p, r := range dataset {
//		fmt.Printf("%q: %q,\n\n", p, r)
//	}
//
//	fmt.Println("Total prompts:", len(dataset))
//}
