package main //tells compiler that is executable file

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

// employees structure will create employees table in database when migrate

type Employees struct {
	gorm.Model        //gorm.model is create id, createdAt, deletedAt and UpdateAt
	EmpName    string `gorm:"emp_name"` //table columns
	EmpPass    string `gorm:"emp_pass"`
	EmpAge     string `gorm:"emp_age"`
	EmpAdd     string `gorm:"emp_add"'`
}

var DB *gorm.DB //global variable for database connection

// database connection function

func DatabaseC(url string) (err error) {
	DB, err = gorm.Open(postgres.Open(url), &gorm.Config{}) //gorm is package and open is function into gorm package and postgres is driver
	return
}

// AddEmp function is used for add employees into the table
func AddEmp(c *gin.Context) {
	var emp Employees

	empname := c.Request.FormValue("emp_name")
	emppass := c.Request.FormValue("emp_pass")
	empage := c.Request.FormValue("emp_age")
	empadd := c.Request.FormValue("emp_add")

	if empname == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("employee name is required"))
		return
	}

	if emppass == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("password is required"))
		return
	}

	if empadd == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("employee address is required"))
		return
	}

	emp.EmpName = empname
	emp.EmpPass = emppass
	emp.EmpAge = empage
	emp.EmpAdd = empadd

	result := DB.Create(&emp) //query for insert the data into database

	if result.Error != nil {
		c.AbortWithError(http.StatusInternalServerError, result.Error)
		return
	}
	c.JSON(http.StatusOK, emp)
	return
}

// GetAllEmp function is used for fetch all employees from the table
func GetAllEmp(c *gin.Context) {
	var emp []Employees
	fmt.Println(emp)
	result := DB.Model(Employees{}).Find(&emp)

	if result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}
	c.JSON(http.StatusOK, emp)
	return
}

// GetOneEmp function is used for fetch One employee from the table where id = id

func GetOneEmp(c *gin.Context) {
	var emp Employees
	id := c.Request.FormValue("id")
	if id == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("id is required"))
		return
	}
	result := DB.Model(Employees{}).Where("id", id).First(&emp)
	if result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}
	c.JSON(http.StatusOK, &emp)
}

// UpdateEmp function is used for Update One employee data into the table where id = id

func UpdateEmp(c *gin.Context) {
	var emp Employees
	empname := c.Request.FormValue("emp_name")
	emppass := c.Request.FormValue("emp_pass")
	empage := c.Request.FormValue("emp_age")
	empadd := c.Request.FormValue("emp_add")

	if empname == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("employee name is required"))
		return
	}

	if emppass == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("password is required"))
		return
	}

	if empadd == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("employee address is required"))
		return
	}
	id := c.Request.FormValue("id")
	if id == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("id is required"))
		return
	}
	result := DB.Model(Employees{}).Where("id=?", id).First(&emp)

	if result.Error != nil {
		c.AbortWithError(http.StatusInternalServerError, result.Error)
		return
	}

	emp.EmpName = empname
	emp.EmpPass = emppass
	emp.EmpAdd = empadd
	emp.EmpAge = empage

	DB.Save(&emp)
	c.JSON(http.StatusOK, &emp)
	return
}

// DeleteEmp function is used for delete One employee data into the table where id = id

func DeleteEmp(c *gin.Context) {
	var emp Employees
	id := c.Request.FormValue("id")

	if id == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("id is required")) //we can use abortwithstats
		return
	}

	result := DB.Model(Employees{}).Where("id", id).First(&emp)
	if result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}
	DB.Delete(&emp)
	c.JSON(http.StatusOK, emp)
}

// program is always start execution with main function
func main() {
	url := "postgres://postgres:shalu123@localhost:5432/user" //database url link
	err := DatabaseC(url)                                     //database connection function
	//error check for database url
	if err != nil {
		fmt.Println("Database problem")
		panic("error in db")
		return
	}
	err = DB.AutoMigrate(Employees{}) // for automatically migrate the table
	//error check for migration is done or not
	if err != nil {
		fmt.Println("error in migration", err)
		panic("error in migration")
		return
	}

	router := gin.Default()
	router.POST("/addemp", AddEmp)      //for add employee details
	router.GET("/allemp", GetAllEmp)    //for see all employees data
	router.GET("/oneemp", GetOneEmp)    //for see only one employee data
	router.PUT("/upemp", UpdateEmp)     //for update only one employee detail
	router.DELETE("/delemp", DeleteEmp) //for delete employee data
	router.Run(":3000")                 //define the port number where it will run

}
