import requests
from xml.etree import ElementTree


# Documentation: 
# https://docs.rice.edu/confluence/display/~lpb1/Courses+API

# For all courses in XML format
# https://courses.rice.edu/courses/!swkscat.cat?format=XML&p_action=COURSE&p_term=201810

# Fall Term   => year + '10'
# Spring Term => year + '20'
# Summer Term => year + '30'
# For fall term add 1 to the year (i.e. 202010 gets fall 2019)


# https://courses.rice.edu/courses/!swkscat.cat?format=XML&p_term=201810&p_action=COURSE
url = "https://courses.rice.edu/courses/!swkscat.cat?format=XML&p_action=COURSE&p_term="
courses = requests.get(url + "202010")
with open("courses.xml", 'w') as file:
    file.write(courses.text)