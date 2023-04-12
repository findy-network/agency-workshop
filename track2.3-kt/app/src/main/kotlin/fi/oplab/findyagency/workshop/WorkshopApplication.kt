package fi.oplab.findyagency.workshop

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RestController

@SpringBootApplication class WorkshopApplication

fun main(args: Array<String>) {
  runApplication<WorkshopApplication>(*args)
}

@RestController
class AppController {

  @GetMapping("/") fun index(): String = "Kotlin example"
  @GetMapping("/greet") fun greet(): String = "IMPLEMENT ME"
  @GetMapping("/issue") fun issue(): String = "IMPLEMENT ME"
  @GetMapping("/verify") fun verify(): String = "IMPLEMENT ME"
}
