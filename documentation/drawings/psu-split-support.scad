$fn=60;
// Power-board backplate
cube([89,7,5]);
translate([3.5,3.5,0]) {
  cylinder(h=10, r=1.25);
  translate([89-7,0,0]) {
    cylinder(h=10, r=1.25);
  }
}
translate([0,70-7,0]){
  cube([89,7,5]);
  translate([3.5,3.5,0]) {
    cylinder(h=10, r=1.25);
    translate([89-7,0,0]) {
      cylinder(h=10, r=1.25);
    }
  }
}
cube([7,70,5]);
translate([89-7,0,0]) 
cube([7,70,5]);

