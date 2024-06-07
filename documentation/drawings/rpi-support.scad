$fn=60;
// RPI backplate
cube([56,7,5]);
translate([3.5,3.5,0]) {
  cylinder(h=10, r=1.25);
  translate([49,0,0]) {
    cylinder(h=10, r=1.25);
  }
}
translate([0,78,0]){
  cube([56,7,5]);
}
translate([0,58,0]) {
  translate([3.5,3.5,0]) {
    cylinder(h=10, r=1.25);
    translate([49,0,0]) {
      cylinder(h=10, r=1.25);
    }
  }
}

cube([7,85,5]);
translate([49,0,0]) 
cube([7,85,5]);

