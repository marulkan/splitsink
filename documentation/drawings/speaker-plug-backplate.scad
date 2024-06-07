
module banana_plug(){
  cylinder(h = 2, r = 2.5, center = false);
  translate([0,20,0])
  cylinder(h = 2, r = 2.5, center = false);
}

module row(){
    translate([0,20,0])
    banana_plug();
    translate([0,114-20-30,0]) 
    banana_plug();
}

difference(){
  cube([146,114,2]);
  
  translate([12,0,0]) {
    row();
  }
  translate([32,0,0]) {
    row();
  }
  translate([53,0,0]) {
    row();
  }
  translate([73,0,0]) {
    row();
  }
  translate([93,0,0]) {
    row();
  }
  translate([113,0,0]) {
    row();
  }
  translate([133,0,0]) {
    row();
  }
}
